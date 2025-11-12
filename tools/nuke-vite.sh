#!/bin/bash

# from the app root (frontend/)
cd frontend || exit 2
rm -rf dist tmp
rm -rf node_modules/.vite
rm -rf .embroider

