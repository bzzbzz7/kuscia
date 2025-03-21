# Copyright 2023 Ant Group Co., Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Configuration file for the Sphinx documentation builder.
#
# This file only contains a selection of the most common options. For a full
# list see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Path setup --------------------------------------------------------------

# If extensions (or modules to document with autodoc) are in another directory,
# add these directories to sys.path here. If the directory is relative to the
# documentation root, use os.path.abspath to make it absolute, like shown here.
#
# import os
# import sys
# sys.path.insert(0, os.path.abspath('.'))

import os.path

# -- Project information -----------------------------------------------------

project = 'Kuscia'
copyright = '2023 Ant Group Co., Ltd.'
author = 'Kuscia authors'

# -- General configuration ---------------------------------------------------

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
extensions = [
    'sphinx.ext.napoleon',
    'sphinx.ext.autodoc',
    'sphinx.ext.graphviz',
    'sphinx.ext.todo',
    'sphinx.ext.viewcode',
    'sphinx.ext.extlinks',
    'sphinx.ext.autosectionlabel',
    'myst_parser',
    "nbsphinx",
    'sphinxcontrib.actdiag',
    'sphinxcontrib.blockdiag',
    'sphinxcontrib.nwdiag',
    'sphinxcontrib.packetdiag',
    'sphinxcontrib.rackdiag',
    'sphinxcontrib.seqdiag',
    'sphinx_design',
]

language = 'zh_CN'
locale_dirs = ['./locales/']  # path is example but recommended.
gettext_compact = False  # optional.

nbsphinx_requirejs_path = ''

# Make sure the target is unique
autosectionlabel_prefix_document = True

# Add any paths that contain templates here, relative to this directory.
templates_path = ['_templates']

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
# This pattern also affects html_static_path and html_extra_path.
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# Enable TODO
todo_include_todos = True

# -- Options for HTML output -------------------------------------------------

# The theme to use for HTML and HTML Help pages.  See the documentation for
# a list of builtin themes.
#
html_theme = 'pydata_sphinx_theme'

# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ['_static']

html_favicon = '_static/favicon.ico'

html_css_files = [
    'css/custom.css',
]

html_js_files = ['js/custom.js']

html_theme_options = {
    "icon_links": [
        {
            "name": "GitHub",
            "url": "https://github.com/secretflow/kuscia",
            "icon": "fab fa-github-square",
            "type": "fontawesome",
        },
    ],
    "logo": {
        "text": "Kuscia",
    },
}

myst_enable_extensions = [
    "amsmath",
    "colon_fence",
    "deflist",
    "dollarmath",
    "fieldlist",
    "html_admonition",
    "html_image",
    "linkify",
    "replacements",
    "smartquotes",
    "strikethrough",
    "substitution",
    "tasklist",
    "attrs_inline",
    "attrs_block",
]

suppress_warnings = ["autosectionlabel", "myst.header"]

# Enable all MyST features
myst_gfm_only = False
# Enable anchors for heading level h1 through h6
myst_heading_anchors = 6


# app setup hook
def setup(app):
    app.add_config_value(
        'recommonmark_config',
        {
            'auto_toc_tree_section': 'Contents',
        },
        True,
    )
