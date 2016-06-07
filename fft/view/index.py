from fft.web import app, template


@app.get("/")
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
