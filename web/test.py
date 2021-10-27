l = ['a', 'ccc', 'aaa', 'ccc', 'aaa']
d = {}
res = {}

for word in l:
    if word in d:
        d[word] += 1
    else:
        d[word] = 1

max = 0
for key in d:
    if d[key] > max:
        max = d[key]
        res.clear()
        res[key] = d[key]
    if d[key] == max:
        res[key] = d[key]

k, v = '', 0
for key, value in res.items():
    if key > k:
        k = key
        v = value
print(k, v)