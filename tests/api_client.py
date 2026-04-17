import requests

class APIClient:
	def __init__(self, base_url):
		self.base_url = base_url
		self.session = requests.Session()

	def _url(self, path):
		return f"{self.base_url}{path}"

	def register(self, login, password):
		return self.session.post(
			self._url("/auth/register"),
			json={
				"login": login,
				"password": password,
			},
		)

	def login(self, login, password):
		return self.session.post(
			self._url("/auth/login"),
			json={
				"login": login,
				"password": password,
			},
		)
