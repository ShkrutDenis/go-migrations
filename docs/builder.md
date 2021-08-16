# Builder API

`import "github.com/ShkrutDenis/go-migrations/builder"`

General, module have a four main function, which using like entrypoint:
 - [NewTable()](#newtable)
 - [ChangeTable()](#changetable)
 - [RenameTable()](#renametable)
 - [DropTable()](#droptable)
 
Module has a next general types:
 - [Table](#table-type)
 - [Column](#column-type)
 - [ForeignKey](#foreign-key-type)
 
#### NewTable

NewTable() will init a structure for adding a new table.

Signature: 

    NewTable(string, *sqlx.DB) Table

Arguments: 
- Name of new table (type string)
- Pointer to DB connection (type *sqlx.DB)

Return value:
- [Table](#table-type)
 
#### ChangeTable

ChangeTable() will init structure for changing an existing table.

Signature: 

    ChangeTable(string, *sqlx.DB) Table

Arguments: 
- Name of existing table (type string)
- Pointer to DB connection (type *sqlx.DB)

Return value:
- [Table](#table-type)
 
#### RenameTable

RenameTable() will init structure for rename an existing table.

Signature: 

    RenameTable(string, string, *sqlx.DB) Table

Arguments: 
- Old name of existing table (type string)
- New name for existing table (type string)
- Pointer to DB connection (type *sqlx.DB)

Return value:
- [Table](#table-type)
 
#### DropTable

DropTable() will init structure for dropping an existing table.

Signature: 

    DropTable(string, *sqlx.DB) Table

Arguments: 
- Name of existing table (type string)
- Pointer to DB connection (type *sqlx.DB)

Return value:
- [Table](#table-type)

## Table type

Table has next methods:
 - [Column()](#column)
 - [String()](#string)
 - [Integer()](#integer)
 - [WithTimeStamps()](#withtimestamps)
 - [RenameColumn()](#renamecolumn)
 - [DropColumn()](#dropcolumn)
 - [PrimaryKey()](#primarykey)
 - [ForeignKey()](#foreignkey)
 - [DropForeignKey()](#dropforeignkey)
 - [GetSQL()](#getsql)
 - [Exec()](#exec)
 - [MustExec()](#mustexec)

#### Column

Column() will init a structure for adding a new column to the table.

Signature: 

    Column(string) Column

Arguments: 
- Name of new column (type string)

Return value:
- [Column](#column-type)

#### String

String() will init a structure for adding a new column with type `varchar(length)` to the table.

Signature: 

    String(string, int) Column

Arguments: 
- Name of new column (type string)
- Length of string (type int)

Return value:
- [Column](#column-type)

#### Integer

Integer() will init a structure for adding a new column with type `int` to the table.

Signature: 

    Integer(string) Column

Arguments: 
- Name of new column (type string)

Return value:
- [Column](#column-type)

#### WithTimeStamps

WithTimeStamps() will add columns `created_at` and `updated_at` to the table.

Signature: 

    WithTimestamps() Table

Return value:
- [Table](#table-type)

#### RenameColumn

RenameColumn() will rename an existing column in the table.

Signature: 

    RenameColumn(string, string) Column

Arguments: 
- Old name of an existing column (type string)
- New name for an existing column (type string)

Return value:
- [Column](#column-type)

#### DropColumn

DropColumn() will drop an existing column from the table.

Signature: 

    DropColumn(string) Column

Arguments: 
- Name of an existing column (type string)

Return value:
- [Column](#column-type)

#### PrimaryKey

PrimaryKey() will add primary key for an existing column or add the new column with type `int` and `autoincrement` to the table.

Signature: 

    PrimaryKey(string) Column

Arguments: 
- Name of an existing or new column (type string)

Return value:
- [Column](#column-type)

#### ForeignKey

ForeignKey() will init structure for adding a new foreign key to the table.

Signature: 

    ForeignKey(string) ForeignKey

Arguments: 
- Name of column in the current table (type string)

Return value:
- [ForeignKey](#foreign-key-type)

#### DropForeignKey

DropForeignKey() will drop foreign key from the table.

Signature: 

    DropForeignKey(string) Table

Arguments: 
- Name of an existing foreign kay (type string)

Return value:
- [Table](#table-type)

#### GetSQL

GetSQL() will return generated sql.

Signature: 

    GetSQL() string

Return value:
- string

*Warning*: result can be a several queries, each query ended on `;`.

#### Exec

Exec() will execute generated sql.

Signature: 

    Exec() error

Return value:
- error

#### MustExec

DropForeignKey() will execute generated sql.

Signature: 

    MustExec()

*Warning*: will be throw `panic()` when error.

## Column type

Column has next methods:
 - [Type()](#type)
 - [Nullable()](#nullable)
 - [NotNull()](#notnull)
 - [Autoincrement()](#autoincrement)
 - [NotAutoincrement()](#notautoincrement)
 - [Default()](#default)
 - [Primary()](#primary)
 - [Unique()](#unique)
 - [NotUnique()](#notunique)
 - [Drop()](#drop)
 - [Change()](#change)
 - [First()](#first)
 - [After()](#after)
 - [Rename()](#rename)

#### Type

Type() will set a type for the column.

Signature: 

    Type(string) Column

Arguments: 
- Valid type for SQL (type string)

Return value:
- [Column](#column-type)

#### Nullable

Nullable() will set `NULL` modifier for the column.

Signature: 

    Nullable() Column

Return value:
- [Column](#column-type)

*Warning*: by default, the column will be created with `NOT NULL` modifier.

#### NotNull
 
NotNull() will set `NOT NULL` modifier for the column.

Signature: 

    NotNull() Column

Return value:
- [Column](#column-type)

#### Autoincrement

Autoincrement() will set `auto_increment` modifier for the column.

Signature: 

    Autoincrement() Column

Return value:
- [Column](#column-type)
 
*Warning*: this method used with [Primary()](#primary) method usually (if you don`t use [PrimaryKey()](#primarykey) method)
*Warning*: no effect for PostgresSQL.

#### NotAutoincrement

NotAutoincrement() will remove `auto_increment` modifier for the column.

Signature: 

    NotAutoincrement() Column

Return value:
- [Column](#column-type)
 
*Warning*: no effect for PostgresSQL. 

#### Default

Default() will set a default value for the column.

Signature: 

    Default(string) Column

Arguments: 
- Valid value for column type (type string)

Return value:
- [Column](#column-type)

#### Primary

Primary() will mark the column as primary key.

Signature: 

    Primary() Column

Return value:
- [Column](#column-type)

#### Unique

Unique() will add for the column a unique key.

Signature: 

    Unique() Column

Return value:
- [Column](#column-type)

#### NotUnique

NotUnique() will remove from the column a unique key.

Signature: 

    NotUnique() Column

Return value:
- [Column](#column-type)

#### Drop

Drop() will mark the column for deleting, then this column will be dropped from table.

Signature: 

    Drop() Column

Return value:
- [Column](#column-type)

#### Change

Change() will mark a column as existed, then column will be change by new modifiers.

Signature: 

    Change() Column

Return value:
- [Column](#column-type)

#### First

First() will set column position as first.

Signature: 

    First() Column

Return value:
- [Column](#column-type)
 
*Warning*: no effect for PostgresSQL.

#### After

After() will set column position after target column.

Signature: 

    After(string) Column

Arguments: 
- Name of existing column (type string)

Return value:
- [Column](#column-type)
 
*Warning*: no effect for PostgresSQL.

#### Rename

Rename() will set a type for the column.

Signature: 

    Rename(string) Column

Arguments: 
- New name for existing column (type string)

Return value:
- [Column](#column-type)

## Foreign Key type

ForeignKey has next methods:
 - [Reference()](#reference)
 - [On()](#on)
 - [OnUpdate()](#onupdate)
 - [OnDelete()](#ondelete)
 - [SetKeyName()](#setkeyname)
 - [GenerateKeyName()](#generatekeyname)
 
#### Reference

Reference() will set a target table name.

Signature: 

    Reference(string) ForeignKey

Arguments: 
- Name of an existing table (type string)

Return value:
- [ForeignKey](#foreign-key-type)

#### On

On() will set a column name related to target table which was set by [Referance()](#reference).

Signature: 

    On(string) ForeignKey

Arguments: 
- Name of an existing column in target table (type string)

Return value:
- [ForeignKey](#foreign-key-type)

#### OnUpdate

OnUpdate() will set a reference option on the update.

Signature: 

    OnUpdate(string) ForeignKey

Arguments: 
- Valid reference option (type string)

Return value:
- [ForeignKey](#foreign-key-type)

#### OnDelete

OnDelete() will set a reference option on delete.

Signature: 

    OnDelete(string) ForeignKey

Arguments: 
- Valid reference option (type string)

Return value:
- [ForeignKey](#foreign-key-type)

#### SetKeyName

SetKeyName() will set a custom name for the foreign kay.

Signature: 

    SetKeyName(string) ForeignKey

Arguments: 
- Name for foreign key (type string)

Return value:
- [ForeignKey](#foreign-key-type)

*Warning*: that method not require. if key will be empty, key name will generate automatically by method [GenerateKeyName()](#generatekeyname).

#### GenerateKeyName

GenerateKeyName() will generate a name for foreign key in the next format: `<base_table>_<target_table>_<target_column>_fk`.

Signature: 

    GenerateKeyName() ForeignKey

Return value:
- [ForeignKey](#foreign-key-type)

*Warning*: that method will be use automatically if foreign key name will be empty.
