#!/usr/bin/env python
from setuptools import setup, find_packages

with open("README.md", "r") as f:
    long_description = f.read()

setup(
    name="tap-incident",
    version="0.5.0",
    description="Singer.io tap for extracting data from the incident.io API",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="incident.io",
    author_email="support@incident.io",
    url="https://github.com/Bilanc/bilanc-incident-io-tap",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
    install_requires=[
        "singer-python>=5.12.1",
        "requests>=2.25.0",
        "python-dateutil>=2.8.2",
        "backoff>=1.11.1",
    ],
    packages=find_packages(),
    package_data={
        "tap_incident": ["schemas/*.json"],
    },
    entry_points={
        "console_scripts": ["tap-incident=tap_incident.cli:main"],
    },
)