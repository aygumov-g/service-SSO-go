from testcontainers.postgres import PostgresContainer
from testcontainers.core.container import DockerContainer
from testcontainers.core.network import Network

from api_client import APIClient
from account_flow import AccountFlow

import requests
import psycopg2
import pytest
import time
import uuid
import os

@pytest.fixture(scope="session")
def network():
	with Network() as net:
		yield net

@pytest.fixture(scope="session")
def postgres(network):
	"""
	Поднимаем чистую PostgreSQL для всех тестов.
	"""

	container = (
		PostgresContainer("postgres:17-alpine")
		.with_network(network)
		.with_network_aliases("postgres")
	)

	with container as pg:
		yield pg

@pytest.fixture(scope="session")
def migrated_for_postgres(network, postgres):
	"""
	Запускаем migrate/migrate контейнер и применяем миграции.
	"""

	timeout = 20
	start = time.time()
	last_error = None
	while time.time() - start < timeout:
		try:
			conn = psycopg2.connect(
				host=postgres.get_container_host_ip(),
				port=postgres.get_exposed_port(5432),
				user=postgres.username,
				password=postgres.password,
				dbname=postgres.dbname,
			)
			conn.close()
			break
		except Exception as e:
			last_error = e
		
		time.sleep(0.5)

	if last_error is not None:
		raise Exception(
			f"Postgres not ready after {timeout}s\n"
			f"Last error: {last_error}\n"
		)

	host_path = os.getenv("HOST_PROJECT_PATH", os.path.abspath("."))
	migrations_path = os.path.join(host_path, "migrations")
	container = (
		DockerContainer("migrate/migrate")
		.with_network(network)
		.with_volume_mapping(migrations_path, "/migrations")
		.with_command([
			"-path", "/migrations",
			"-database", "postgres://%(user)s:%(pass)s@postgres:5432/%(db)s?sslmode=disable" % {
				"user": postgres.username,
				"pass": postgres.password,
				"db": postgres.dbname,
			},
			"up",
		])
	)

	with container as ct:
		result = ct.get_wrapped_container().wait()
		logs = ct.get_wrapped_container().logs().decode()
		if result["StatusCode"] != 0:
			raise Exception(
				"Migration failed\n"
				f"Container logs:\n{logs}"
			)

	return postgres

@pytest.fixture(scope="session")
def app(network, migrated_for_postgres):
	"""
	Запускаем контейнер с копией основного go-сервиса. 
	"""

	container = (
		DockerContainer("aygumov-g/service-sso-app-go")
		.with_network(network)
		.with_exposed_ports(8080)
		.with_env("APP_PORT", "8080")
		.with_env("POSTGRES_HOST", "postgres")
		.with_env("POSTGRES_USER", migrated_for_postgres.username)
		.with_env("POSTGRES_PASSWORD", migrated_for_postgres.password)
		.with_env("POSTGRES_DB", migrated_for_postgres.dbname)
		.with_env("JWT_SECRET", "secret")
		.with_env("JWT_TTL", "5m")
		.with_env("REFRESH_TTL", "170h")
	)

	with container as ct:
		timeout = 20
		start = time.time()
		last_error = None
		while time.time() - start < timeout:
			try:
				if requests.get(f"http://{ct.get_container_host_ip()}:{ct.get_exposed_port(8080)}/health/ready").status_code == 200:
					last_error = None
					break
			except Exception as e:
				last_error = e
			
			time.sleep(0.5)
		
		logs = ct.get_wrapped_container().logs().decode()
		if last_error is not None:
			raise Exception(
				f"App not ready after {timeout}s\n"
				f"Last error: {last_error}\n\n"
				f"Container logs:\n{logs}"
			)

		yield ct

@pytest.fixture()
def account(app):
	return AccountFlow(
		APIClient(f"http://{app.get_container_host_ip()}:{app.get_exposed_port(8080)}")
	)
