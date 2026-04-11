import os
import random
import string
import pytest
import requests

class APIClient:
	def __init__(self, base_url):
		self.base_url = base_url
		self.session = requests.Session()

	def post(self, path, **kwargs):
		return self.session.post(self.base_url + path, **kwargs)

	def get(self, path, **kwargs):
		return self.session.post(self.base_url + path, **kwargs)

@pytest.fixture
def client():
	return APIClient("http://%(host)s:%(port)s" % {
		"host": os.environ["APP_HOST"],
		"port": os.environ["APP_PORT"],
	})

@pytest.fixture
def test_user():
	return {
		"login": "".join([random.choice(string.ascii_letters + string.digits) for i in range(9)]),
		"password": "12345",
	}
