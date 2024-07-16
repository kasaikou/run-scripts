# Markflow サンプル

## シェルスクリプトの実行

シェルで `echo "hello world"` を実行します。

```yaml
kind: execute
content:
  environments:
    - TEST=environment
```

```sh
echo "hello world"
```

## Python の実行

Python で `print("hello world")` を実行します。

```py
print("hello world")
```
