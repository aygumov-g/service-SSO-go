def test_register_successful(account):
	account.register().status_code_equals(201)

def test_register_conflict_login(account):
	account.register().register().status_code_equals(409)
