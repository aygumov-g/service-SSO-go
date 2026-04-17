import uuid

class AccountFlow:
	def __init__(self, api):
		self.api = api

		"""
		Тип uuid в логине всего-лишь ради примера, на практике именно он - не требуется.
		Ключ "fake_password" вообще необязателен и не должен присутствовать в объекте
		представления данных об аккаунте.
		"""
		self.account_data = {
			"login": str(uuid.uuid4()),
			"password": "12345",
			"fake_password": "123",
		}

		"""
		Любой запрос с использованием данного класса - напрямую
		влияет на содержимое переменной ниже.
		"""
		self.response = None
	
	def register(self):
		self.response = self.api.register(
			login=self.account_data["login"],
			password=self.account_data["password"],
		)

		return self

	def login(self, is_fake_password=False):
		self.response = self.api.login(
			login=self.account_data["login"],
			password=self.account_data["fake_password"] if is_fake_password else self.account_data["password"],
		)

		return self
	
	def status_code_equals(self, code):
		assert self.response != None and self.response.status_code == code, (
			f"Code {code}, got {self.response.status_code}, "
			f"body={self.response.text}"
		)

		return self
