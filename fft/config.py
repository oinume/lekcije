# coding: utf-8
import os

config_common = {
    "debug": False,
    "host": "0.0.0.0",
    "port": 5000,
    "reloader": False,
    "server": "gunicorn",
}
config_local = {
    "debug": True,
    "reloader": True,
    # "server": "wsgiref",
    #    "server": "gunicorn",
    #    "workers": 1,
}
config_production = {
    "debug": False,
    "reloader": False,
    "workers": 1,
}

mode = os.environ.get("DMM_EIKAIWA_FFT_ENV", "local")
config_env = None
if mode == "local":
    config_env = config_local
elif mode == "production":
    config_env = config_production
else:
    raise ValueError("Unknown config mode: " + mode)

config = config_common.copy()
config.update(config_env)
