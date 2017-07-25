Configuration
=============

`.mailfull/config`

| key           | type   | default                                   | required | description                                                     |
|:--------------|:-------|:------------------------------------------|:---------|:----------------------------------------------------------------|
| dir_database  | string | `"./etc"`                                 | no       | A relative path from repository dir (or a absolute path)        |
| dir_maildata  | string | `"./domains"`                             | no       | A relative path from repository dir (or a absolute path)        |
| username      | string | The username who executed `mailfull init` | **yes**  | It used for setting owner of database files and maildata files. |
| cmd_postalias | string | `"postalias"`                             | no       | Command name or path                                            |
| cmd_postmap   | string | `"postmap"`                               | no       | Command name or path                                            |
