from setuptools import setup

NAME = 'mongoke'
setup(
    name=NAME,
    # [bump]
    version='0.0.15',
    description=NAME,
    long_description=open('README.md').read(),
    long_description_content_type='text/markdown',

    author='Tommaso De Rossi',
    author_email='daer.tommy@gmail.com',
    license='Apache Software License 2.0',
    url=f'https://github.com/remorses/{NAME}',
    keywords=['TODO'],
    install_requires=[x for x in open('./requirements.txt').read().strip().split('\n') if x.strip()],
    package_data={'': ['*.yaml', '*.json', '*.yml', 'VERSION', 'README.md', 'requirements.txt']},
    classifiers=[
        'Intended Audience :: Information Technology',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
    ],
    packages=['mongoke', 'mongoke.templates', ],
    scripts=[
        'bin/mongoke'
    ]
)


