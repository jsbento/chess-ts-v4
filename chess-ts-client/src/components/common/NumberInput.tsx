import React, { ChangeEvent } from 'react'

interface NumberInputProps {
  label?: string
  labelClassName?: string
  inputClassName?: string
  value: number
  onChange: (value: number) => void
}

const NumberInput: React.FC<NumberInputProps> = ({
  label,
  labelClassName,
  inputClassName,
  value,
  onChange,
}) => {
  return (
    <label className={labelClassName}>
      {label}
      <input
        type='number'
        className={`ml-2 w-20 px-2 text-right ${inputClassName}`}
        value={value}
        onChange={(e: ChangeEvent<HTMLInputElement>) => {
          const value = parseInt(e.target.value)
          onChange(isNaN(value) ? 0 : value)
        }}
      />
    </label>
  )
}

export default NumberInput
