package decoding

const PQ_TABLE = `(
    (
        (element
            (text) @text) 
        .
        (element
            (start_tag
                (tag_name) @tn)) @el
    )
    (#eq? @text "~text")
    (#eq? @tn "table")
)`
