# import the Flask instance named flask_app from the recipe_app package
from recipe_app import flask_app
from flask import render_template


# Route to home webpage
@flask_app.route("/")
@flask_app.route("/index", methods=["GET", "POST"])
def index():
    return render_template("index.html")


# Route to example page
@flask_app.route("/example")
def example():
    return render_template('example_temp.html')
