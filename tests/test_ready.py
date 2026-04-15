def test_dummy(postgres):
	assert postgres is not None

def test_migrations(migrated_for_postgres):
	assert migrated_for_postgres is not None

def test_app_started(app):
	assert app is not None

def test_ready_successful(app):
	assert __import__("requests").get(f"http://{app.get_container_host_ip()}:{app.get_exposed_port(8080)}/health/ready").status_code == 200
