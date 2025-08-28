from setuptools import setup

setup(
    name='codelearn',
    version='1.0.0',
    py_modules=['codelearn'],
    install_requires=['requests'],
    entry_points={
        'console_scripts': [
            'codelearn=codelearn:main',
        ],
    },
)
