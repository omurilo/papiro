# Papiro

Papiro é um gerador de sites estáticos super rápido e minimalista, escrito em Go. Ele transforma arquivos Markdown em HTML de forma simples e eficiente, focando em produtividade e facilidade de uso.

## Características

- **Conversão rápida** de arquivos Markdown para HTML
- **Foco em simplicidade e minimalismo**
- **Templates personalizáveis**
- **Ideal para blogs, portfólios e sites pessoais**

## Instalação

Você pode compilar o Papiro a partir do código-fonte:

```bash
git clone https://github.com/seu-usuario/papiro.git
cd papiro
go build -o papiro ./cmd/papiro/
```

## Uso

No terminal, execute:

```bash
./papiro
```

Veja os comandos disponíveis com:

```bash
./papiro --help
```

## Exemplo de uso

### Comando `init`

Para começar rapidamente um novo projeto com a estrutura básica, use:

```bash
./papiro init
```

Esse comando cria as pastas essenciais (`content` e `theme/static`), copia os templates padrão para a pasta `theme` e gera um exemplo de post Markdown em `content/hello-world.md`. Assim, você pode começar a escrever imediatamente sem se preocupar com a configuração inicial.

### Comando `build`

Para gerar um site a partir de arquivos Markdown:

```bash
./papiro build
```

- A pasta `examples` trás um exemplo completo de como pode ser um projeto papiro, a pasta `examples/first-blog/theme` é completamente opcional.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues e pull requests.

## Licença

Este projeto está licenciado