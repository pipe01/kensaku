grammar Query;

query
    : .*
    ;

quoted
    : '"' .*? '"'
    ;

operator
    : '(' name=.+? ':' value=.+? ')'    # opField
    | quoted                            # opExact
    | .+?                               # opText
    ;
