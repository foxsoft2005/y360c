y360c is a command line interface (cli) for Yandex 360, which uses the standard API to access the corresponding entities/actions in Yandex 360.

There is no need to install the application, just place the exe file anywhere and start working with it, like with any other console application.
If necessary, you can add the path to the application to PATH.

Before starting work, you will need to create an application with the necessary access rights in https://oauth.yandex.ru/, as well as obtain an auth token, which is necessary to access data in Yandex360.

The token can be used in 2 ways:
1. pass it explicitly in command using the flag `--token (-t)`
2. save the token for automatic use in the application configuration file, which is located here: `%USERPROFILE%\y360c.json` (created automatically when using the command `init` - i.e. `y360c init`).

In addition to the token, most operations require the "Organization Id" parameter, which must be used similarly to the token:
1. can be passed explicitly in the parameters using the flag `--orgId (-o)`
2. can be saved for automatic use in the application configuration file.

Examples:

**List of organizations**:

`y360c org ls` (if the token is specified in the config)

or

`y360c org ls --token <token>`

**List of departments**:

`y360c dept ls --orgId <organizaion id> --token <token>`

or

`y360c dept ls` (if the token and the organization id are specified in the config)

**List of employees**:

`y360c user ls` (if the token and the organization id are specified in the config)

or

`y360c user ls --orgId <organizaion id> --token <token>`

or (search for a specific employee)

`y360c user ls --orgId <organizaion id> --token <token> --id <employee id>`

or (search my email)

`y360c user ls --orgId <organizaion id> --token <token> --email aaa@vvv.cc`

or (search by first/middle/last name)

`y360c user ls --orgId <organizaion id> --token <token> --name vladimir`
