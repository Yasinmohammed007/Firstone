from concurrent import futures
import os,json
import requests
import pytest

def test_success():    
    response = requests.get("http://config-service:8086/test")
    # response_body = response.json()
    # print("test",response.status_code)
    assert response.status_code == 200
    assert response.text == "Welcome to the API server!!"
