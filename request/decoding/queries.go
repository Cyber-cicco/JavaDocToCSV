package decoding

const (
    Q_LINK= `(
(start_tag
  (tag_name) @tag_name
  (attribute 
    (attribute_name)
    (quoted_attribute_value 
      (attribute_value) @name))) 
(#eq? @name "typeSummary")
(#eq? @tag_name "table")
)`
)


