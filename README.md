### Deployer

Deployer é o componente de plataforma responsável por localizar instâncias que devem ser reprocessadas caso a persistência de uma instância em execução seja realizada.

### Build

O build da aplicação é feito através de um arquivo Makefile, para buildar a aplicação execute o seguinte comando:

```sh
$ make
```

Após executar o make será criada uma pasta dist e o executável da aplicação deployer.

### Deploy

O processo de deploy do deployer na plataforma é feito através do installer, os componentes em Go são compilados e comitados dentro do installer então para atualizar a versão do deployer para atualizar a versão do deployer na plataforma utilize o seguinte comando:

```sh
$ mv dist/deployer ~/installed_plataforma/Plataforma-Installer/Dockerfiles
$ plataforma --upgrade deployer
```

### API

Criar uma nova solution

```http
POST /api/v1.0.0/solution HTTP/1.1
Host: localhost:6970
Content-Type: application/json

{
  "name":"nome",
  "version":"1313123123",
  "id":"a2806dce-e84d-4359-952b-1a514ae74fac",
  "description":"minha solution"
}
```

Adicionar chave publica

```http
POST /api/v1.0.0/publickey/<solution_name>/<name> HTTP/1.1
Host: localhost:6970
Cache-Control: no-cache

<sua chave publica>
```

Fazer o deploy de uma app

```http
POST /api/v1.0.0/app/<app_id>/deploy HTTP/1.1
Host: localhost:6970
Content-Type: application/json

```

Criar uma nova app

```http
POST /api/v1.0.0/solution/<solution_id>/create/app HTTP/1.1
Host: localhost:6970
Content-Type: application/json

{
  "name":"app",
  "version":"1.0.0",
  "id":"b2806dce-e84d-4359-952b-1a514ae74fac",
  "description":"minha app",
  "type":"process"
}
```


### Organização do código

1. actions
    * São as principais ações do serviço, por exemplo, criar apps e fazer deploy;
2. api
    * É a declaração da API do deployer;
3. container
    * É um pacote para montar um container docker
4. env
    * Pacote de funções para lidar com variável de ambiente
5. models
    * Define o modelo de domínio usado pelo deployer
6. git
    * Pacota com algumas funcionalidades do git como, por exemplo, o clone
7. sdk
    * Implementa chamadas dos serviços da plataforma
8. vendor
    * É um pacote do Go onde ficam todas as bibliotecas de terceiros, os arquivos deste pacote jamais devem ser alterados diretamente;