horsebase
======================

**horsebase** creates a database for horse racing analysis.  

# Usage
`horsebase` stores the horse racing data in the database by executing the following command.  
`  
$ horsebase -build  
`  

# install
Build by yourself.  
`  
$ make install  
`  

Other installation methods(ex. go get) will be supported in the future.  

# Configuration File
`horsebase` provides a configuration file in toml format.
* [file/horsebase.toml](#horsebasetoml)
* [file/bloodtype.toml](#bloodtypetoml)

## horsebase.toml

#### config
You can specify the start date of the data to be stored.

`  
[config]  
~~~  
oldest_date=20070101  
~~~  
`  

#### db
You can set the database username and password.  

`  
[db]  
dbuser = "$username"  
dbpass = "$password"  
`

## bloodtype.toml
You can customize the pedigree information and register it in the database.  

You can define the lineage.  

`  
[bloodtype]  
bloodtypes = [  
  'AAA系',  
  ~~~  
]  
`  

You can define main blood-type and sub blood-type and you can map stallions to it.

`  
[mainbloodtypes]  

  [mainbloodtypes.'AAA系']  
  stallions = [  
    'sample stallion A1',  
    'sample stallion A2',  
    ~~~
  ]  

  [mainbloodtypes.'BBB系']  
    stallions = [  
      'sample stallion B1'  
      'sample stallion B2',  
      ~~~  
    ]  

  ~~~~  

  [subbloodtypes.'CCC系']  
  stallions = [  
    'sample stallion C1',  
    'sample stallion C2',  
    ~~~  
  ]  

  ~~~~  

`

`horsebase` provides templates and you can use it.

# Command Line Options
`  

--build            Stores all data  

--init_db          Create horsebase DB  
--reg_bloodtype    Store the bloodtype data defined in bloodtype.toml in horsebase DB  
--make_list        Save the URL of the race data in racelist.txt  
--get_racehtml     Gets the HTML form the URL listed in racelist.txt  
--reg_racedata     Scrape HTML and store race data in horsebase DB  
--reg_horsedata    Scrape HTML and store horse data in horsebase DB  
--drop_db          Delete horsebase DB  
--match_bloodtype  Map bloodtype data and stallion data defined in bloodtype.toml  
  
`  
