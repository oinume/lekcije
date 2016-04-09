#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
from sys import argv
import bottle
from bottle import template, response, get
#from tsc.models import DB

bottle.debug(True)
application = bottle.default_app()
if __name__ == "__main__":
    bottle.run(application, host="0.0.0.0", port=argv[1])


@get("/")
def index():
    html = """
<html>
<head>
<meta charset="UTF-8"></meta>
<title>dmm-eikaiwa-fft - Follow favorite teachers of DMM Eikaiwa</title>
</head>
<body>
hello
</body>
</html>
    """
    return template(html)


# @get("/status")
# def status():
#     response.content_type = "application/json; charset=utf-8"
#     conn = None
#     try:
#         conn = DB.connect(os.environ.get("CLEARDB_DATABASE_URL"))
#         with conn.cursor() as cursor:
#             cursor.execute("SELECT COUNT(*) FROM teacher")
#             cursor.fetchone()
#         return {
#             "APP_ID": os.environ.get("APP_ID"),
#             "db": "true",
#         }
#     finally:
#         if conn:
#             conn.close()
