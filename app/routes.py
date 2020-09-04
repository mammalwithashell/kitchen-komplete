# import the Flask instance named app from the app package
from app import app
from flask import render_template


# Route to home webpage
@app.route("/")
@app.route("/index", methods=["GET", "POST"])
def index():
    return "Hello, World!"


# Route to example page
@app.route("/example")
def example():
    return render_template('example_temp.html')
