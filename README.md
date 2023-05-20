Español

Decisiones

- Se creo el paquete `commons` con tres propósitos
- Aislar a través de interfaces la dependencia externa (tipo conexión con Redis, con SQS de AWS), intentando tener un nivel bajo de acoplamiento.
- Creación de mocks para poder testear más fácil los demás componentes (mock de la fecha actual).
- Compartir lógica entre las dos lambdas, ejemplo, el modelo (`struct`).
- Se guardó las reglas de límite de notificación en un repositorio (`rate-limit/v1/config_repository.go`) que mañana
  puede hacer refactor para conectarse a una fuente externa, para así no tener que desplegar la lambda para tener
  cambios.

Instrucciones

Ubicar archivo serverless.yml y cambiar las variables de entorno

English

Decisions

- The `commons` package was created for three purposes
- Isolate the external dependency through interfaces (connection type with Redis, with AWS SQS), trying to have a low level of coupling.
- Creation of mocks to be able to test the other components more easily (mock of the current date).
- Share logic between the two lambdas, for example, the model (`struct`).
- Saved notification limit rules to a repository (`rate-limit/v1/config_repository.go`) tomorrow
  you can refactor to connect to an external source, so you don't have to deploy the lambda to have
  changes.

Instructions

Locate serverless.yml file and change environment variables


Envs
```
REDIS_HOST: "127.0.0.1"
REDIS_PORT: "6379"
REDIS_DB: 0
REDIS_PASSWORD: ""
QUEUE_URL: "queue-url"

MAIL_HOST: "smtp.gmail.com"
MAIL_USERNAME: "kennitromero@gmail.com"
MAIL_PASSWORD: "password-is-a-password"
MAIL_PORT: "25"
```
