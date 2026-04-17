import requests

class APIClient:
	def __init__(self, base_url):
		self.base_url = base_url
		self.session = requests.Session()

	def _url(self, path):
		return f"{self.base_url}{path}"

	def register(self, login, password):
		return self.session.post(
			url=self._url("/auth/register"),
			json={
				"login": login,
				"password": password,
			},
		)

	def login(self, login, password):
		return self.session.post(
			url=self._url("/auth/login"),
			json={
				"login": login,
				"password": password,
			},
		)
	
	def me(self, access_token):
		return self.session.get(
			url=self._url("/auth/me"),
			headers={
				"Authorization": f"Bearer {access_token}",
			},
		)

	def change_password(self, access_token, old_password, new_password):
		return self.session.post(
			url=self._url("/auth/change_password"),
			json={
				"old_password": old_password,
				"new_password": new_password,
			},
			headers={
				"Authorization": f"Bearer {access_token}",
			},
		)

	def refresh(self, refresh_token):
		return self.session.post(
			url=self._url("/auth/refresh"),
			json={
				"refresh_token": refresh_token,
			},
		)
