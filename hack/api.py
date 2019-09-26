#!/bin/sh python

#
# To use this example script, run kubectl proxy on your developer machine on port 8001 by running,
# [kubectl proxy --port=8001]
#
import json
import requests
import sys

apiVersion = "logging.pf9.io/v1alpha1"
ownerName = "example-script"

class Api(object):
    def __init__(self):
        self.bUrl = 'http://localhost:8001/apis/' + apiVersion

    def get(self, path):
        url = '/'.join([self.bUrl, path])
        r = requests.get(url=url, headers={"Content-Type": "application/json"})
        r.raise_for_status()
        return r.json()
    
    def post(self, path, data):
        url = '/'.join([self.bUrl, path])
        r = requests.post(url=url, json=data, headers={"Content-Type": "application/json"})
        r.raise_for_status()
    
    def patch(self, path, data):
        url = '/'.join([self.bUrl, path])
        r = requests.patch(url, json=data, headers={"Content-Type": "application/json-patch+json"})
        r.raise_for_status()
    
    def delete(self, path):
        url = '/'.join([self.bUrl, path])
        r = requests.delete(url=url, headers={"Content-Type": "application/json"})
        r.raise_for_status()

def getOutput(api, name):
    print(api.get('outputs'))

def createOutput(api, name):
    data = {
        "kind": "Output",
        "apiVersion": apiVersion,
        "metadata": {
            "name": name,
            "labels": {
				"created_by": ownerName,
			}
        },
        "spec": {
            "type": "elasticsearch",
            "params": [
                {
                    "name": "url",
                    "value": "http://elasticsearch.default.cluster.local:9200"
                },
                {
                    "name": "user",
                    "value": "test-elastic"
                },
                {
                    "name": "password",
                    "value": "test-password"
                },
                {
                    "name": "index_name",
                    "value": "test-index"
                }
            ]
        }
    }
    api.post('outputs', data)

def deleteOutput(api, name):
    api.delete('outputs/' + name)

if __name__ == "__main__":
    if len(sys.argv) != 2 or sys.argv[1] not in ["get", "create", "delete"]:
        sys.stderr.write("Usage %s <get|create|delete>\n" % sys.argv[0])
        exit(1)

    api = Api()
    if sys.argv[1] == "get":
        getOutput(api, "pf9-example")
    elif sys.argv[1] == "create":
        createOutput(api, "pf9-example")
    else:
        deleteOutput(api, "pf9-example")