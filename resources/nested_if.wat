(module
  (func $nested_if (export "nested_if") (param i32) (param i32) (result i32)
    local.get 0
    local.get 1
    i32.lt_s
    
    if (result i32)
      local.get 0
      local.get 1
      i32.lt_s

      if (result i32)
        local.get 0
        i32.const 10
        i32.mul
      else
        unreachable
      end
    else 
      i32.const 3
    end
    unreachable
  )
)
