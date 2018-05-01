import sys
import os
from setuptools import setup

if sys.version_info.major == 3:
    base_dir = 'py3'
else:
    assert sys.version_info.major == 2
    assert sys.version_info.minor >= 7
    base_dir = 'py2'

setup(
    name='dispy',
    version='4.8.6',
    description='Distributed and Parallel Computing with/for Python.',
    long_description=open('README.rst').read(),
    keywords='distributed computing, parallel processing, mapreduce, hadoop, job scheduler',
    url='http://dispy.sourceforge.net',
    author='Giridhar Pemmasani',
    author_email='pgiri@yahoo.com',
    package_dir={'':base_dir},
    packages=['dispy'],
    package_data = {
        'dispy' : ['data/*', 'examples/*', 'doc/*'],
    },
    install_requires=['pycos >= 4.6.5'],
    scripts=[os.path.join(base_dir, 'dispy', script) for script in \
             ['dispynode.py', 'dispynetrelay.py', 'dispyscheduler.py']] + \
            [os.path.join(base_dir, script) for script in ['dispy.py']],
    license='Apache 2.0',
    platforms='any',
    classifiers=[
        'Development Status :: 5 - Production/Stable',
        'Intended Audience :: Developers',
        'Intended Audience :: Science/Research',
        'License :: OSI Approved :: Apache Software License',
        'Operating System :: OS Independent',
        'Programming Language :: Python :: 2',
        'Programming Language :: Python :: 2.7',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.1',
        'Topic :: Scientific/Engineering',
        'Topic :: Software Development',
        ]
    )
