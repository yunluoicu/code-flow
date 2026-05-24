#!/usr/bin/env python3
import sys, subprocess, pathlib
subprocess.call([sys.executable, str(pathlib.Path(__file__).with_name('install_all.py')), '--tools', 'cursor'] + sys.argv[1:])
