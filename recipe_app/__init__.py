"""__init__.py is a constructor document that lets the python interpreter 
    see the folder named app as if it were a package that we can import objects from

    First import the Flask class from the flask library
    __name__ refers to the name in reference to the main python process

    Next import the pages of the web app stored in routes
    """
from flask import Flask
from flask_sslify import SSLify
import flask_sqlalchemy

flask_app = Flask(__name__)
flask_app.debug = True
sslify = SSLify(flask_app)

from recipe_app import routes

# Import routes into the package
