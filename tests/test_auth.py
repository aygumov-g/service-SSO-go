import pytest

def test_register(client, test_user):
    rest = client.post("/auth/register", json=test_user)
    assert rest.status_code == 201
