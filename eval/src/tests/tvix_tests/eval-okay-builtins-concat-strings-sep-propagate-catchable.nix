map (e: (builtins.tryEval e).success) [ (builtins.concatStringsSep (builtins.throw "a") [ "" ]) (builtins.concatStringsSep "," (builtins.throw "a")) (builtins.concatStringsSep "," [ "a" (builtins.throw "a") ]) ]