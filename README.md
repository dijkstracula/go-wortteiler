# Wortteiler

CI status: [![CircleCI](https://circleci.com/gh/dijkstracula/go-wortteiler/tree/master.svg?style=svg)](https://circleci.com/gh/dijkstracula/go-wortteiler/tree/master)

## A golang re-implementation of [Bandwortersplitter](https://github.com/dijkstracula/Bandwortersplitter).

### Translation backend

Sample output:
```
$ curl -s -XPOST 'localhost:8080/split/entschuldigung' | json_pp
{
   "defn" : "Entschuldigung /ɛntʃuldiguŋ/ <n, s>\n apology; excuse; excuse me; sorry\n",
   "prefix" : {
      "defn" : "ent"
   }
   "suffix" : {
      "prefix" : {
         "defn" : "schuldig /ʃuldiç/\n blamable; due; guilty\n",
         "suffix" : {
            "defn" : "ig"
         },
         "prefix" : {
            "defn" : "Schuld /ʃult/ <n, s>\n blame; debt; guilt; guiltiness\n"
         }
      }
      "suffix" : {
         "defn" : "ung"
      },
   },
}
```
