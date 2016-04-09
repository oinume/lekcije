#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
from sys import argv
import bottle
from bottle import static_file, Bottle, TEMPLATE_PATH, template, get
#from tsc.models import DB

bottle.debug(True)
application = bottle.default_app()
app_root = os.path.dirname(os.path.abspath(__file__))
app = Bottle()
app.config['app_root'] = app_root
#app.config.update(config)
TEMPLATE_PATH.append(os.path.join(app_root, 'templates'))

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


@get('/static/<file_name:re:.+>', name='static')
def serve_static(file_name):
    return static_file(file_name, root=os.path.join(app.config['app_root'], 'static'))


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
