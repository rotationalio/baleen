# Baleen Dashboard
This directory contains a work-in-progress frontend for Baleen in the form of an interactive dashboard.

## Running the Dashboard
The current dashboard utilizes pyscript to run Python code in the browser and render interactive visualizations of the corpus data. Currently this relies on data being served locally, so to start the HTTP server there is a helper script that can be run from this directory.

`python mock/server.py`

Once the server is up and running, simply open `vocab.html` in a web browser and wait for the dashboard to load. Note: It may take 10-15 seconds to load the external dependencies and initialize the dashboard.