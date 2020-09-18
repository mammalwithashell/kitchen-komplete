"""__init__.py is a constructor document that lets the python interpreter 
    see the folder named app as if it were a package that we can import objects from

    First import the Flask class from the flask library
    __name__ refers to the name in reference to the main python process

    Next import the pages of the web app stored in routes
    """
from flask import Flask
import flask_sqlalchemy
from config import Config # import config.py
from flask_talisman import Taslisman


flask_app = Flask(__name__)
flask_app.debug = True
flask_app.config.from_object(Config)

talisman = Taslisman(flask_app)

from recipe_app import routes # Import routes into the package
