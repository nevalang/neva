# About

"99 Bottles of Beer" is a classical programming task that utilizes loops, conditions and io. You can see the details [here](https://www.99-bottles-of-beer.net).

## TODO

Is it possible to simplify the program? Example (Python):

```python
for quant in range(99, 0, -1):
   if quant > 1:
      print(quant, "bottles of beer on the wall,", quant, "bottles of beer.")
      if quant > 2:
         suffix = str(quant - 1) + " bottles of beer on the wall."
      else:
         suffix = "1 bottle of beer on the wall."
   elif quant == 1:
      print("1 bottle of beer on the wall, 1 bottle of beer.")
      suffix = "no more beer on the wall!"
   print("Take one down, pass it around,", suffix)
   print("--")
```
