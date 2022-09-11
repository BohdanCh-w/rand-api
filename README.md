# RandAPI
### RandAPI is an cli tool to retrieve data from random.org

---

## Usage
```
randapi [global options] command [command options] [arguments...]
```

---

## Commands

| Name    | Aliases | Description                                           |
| ------- | ------- | ----------------------------------------------------- |
| integer | int     | generate random integer in range (including)          |
| coin    |         | generate random coinflip result (two values possible) |
| decimal | dec     | generate random decimal value in range [0, 1]         |
| gausian | gaus    | generate random value with Gausian distribution       |
| string  | str     | generate random string of given characters            |
| uuid    |         | generate random uuid V4                               |
| blob    |         | generate random Binary Large OBject                   |
| status  | st      | get specified apiKey usage                            |
| help    | h       | show a list of commands or help for one command       |

---

## Global Options

| Name            | Aliases | Description                                               | Default Value                          |
| --------------- | ------- | --------------------------------------------------------- | -------------------------------------- |
| apikey value    |         | specify custom [API-key](https://api.random.org/api-keys) | embeded resource if any, else required |
| file value      | -f      | save output to specied file                               | \<STDOUT\>                             |
| help            | -h      | show help                                                 | false                                  |
| quite           | -q      | suppress all warnings                                     | false                                  |
| separator value | --sep   | string to separate output                                 | " "                                    |
| signed          | -s      | get signed reply from random.org                          | false                                  |
| timeout value   | -t      | randomness server response timeout in seconds             | 5                                      |
| verbose         | -v      | make verbose output after completition                    | false                                  |

To get options for specific command use `randapi [command] -h`

---

## TODO List

- Finish documentation
- Add integer sequenses command
- Add Signed requests
- Add Pregenerated requests
