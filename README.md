dev-up allows you to specify a configuration for your development environment
so that your environment setups are not painfully long and full of outdated information.
dev-up eliminates the need for having work environment setup documentation that you  have to keep updated.

## Yaml configuration example
```yaml
profile: dev-environment-setup
    execute:
        - program: custom-program
        - command: get config
    lifecycle:
        - install:
            app: python
            version: 13.6.5
            env:
                - name: something1
                value: value of something
                - name: something2
                value: something2-value
        - install:
            app: node
            version: 13.6.5
            env:
                - name: something1
                value: value of something
                - name: something2
                value: something2-value
```
