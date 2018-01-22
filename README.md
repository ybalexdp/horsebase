horsebase
======================

**horsebase** creates a database for horse racing analysis.  

# Usage
`horsebase` stores the horse racing data in the database by executing the following command.  
```bash
$ horsebase -b
```  

# install

## Binary released
You can download the binary from [release page](https://github.com/ybalexdp/horsebase/releases).

## Homebrew
```bash
$ brew tap ybalexdp/horsebase
```

```bash
$ brew install horsebase
```

## Build by yourself.  
```bash
$ make install  
```  
## Go Get
Run the following command beforehand.
```bash
$ export PATH=$PATH:$GOPATH/bin
```

You can install by "go get".  
```bash
$ go get github.com/ybalexdp/horsebase  
```

and

```bash
$ cd $GOPATH/src/github.com/ybalexdp/horsebase  
```

# Configuration File
`horsebase` provides a configuration file in toml format.
* [file/horsebase.toml](#horsebasetoml)
* [file/bloodtype.toml](#bloodtypetoml)

## horsebase.toml

### config
You can specify the start date of the data to be stored.

```bash
[config]  
~~~
oldest_date=20070101  
~~~
```

### db
You can set the database username and password.  

```bash
[db]  
dbuser = "$username"  
dbpass = "$password"  
```

## bloodtype.toml
You can customize the pedigree information and register it in the database.  

You can define the lineage.  

```bash
[bloodtype]  
bloodtypes = [  
  'AAA系',  
  ~~~  
]  
```

You can define main blood-type and sub blood-type and you can map stallions to it.

```bash
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

```

`horsebase` provides templates and you can use it.

# Command Line Options
```bash
--build,-b            # Store all data  

--init_db,-i          # Create horsebase DB  
--reg_bloodtype       # Store the bloodtype data defined in bloodtype.toml in horsebase DB  
--list,-l             # Save the URL of the race data in racelist.txt  
--get_racehtml        # Get the HTML form the URL listed in racelist.txt  
--reg_racedata        # Scrape HTML and store race data in horsebase DB  
--reg_horsedata       # Scrape HTML and store horse data in horsebase DB  
--drop_db,-d          # Delete horsebase DB  
--match_bloodtype,-m  # Map bloodtype data and stallion data defined in bloodtype.toml  
--update,-u           # Collect and store recent race data from last stored data

```
