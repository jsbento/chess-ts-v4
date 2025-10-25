import React from 'react'

interface CheckboxProps {
  label?: string
  labelClassName?: string
  inputClassName?: string
  checked: boolean
  onChange: (checked: boolean) => void
}

const Checkbox: React.FC<CheckboxProps> = ({
  label,
  labelClassName,
  inputClassName,
  checked,
  onChange,
}) => {
  return (
    <label className={labelClassName}>
      <input
        type='checkbox'
        className={`mr-2 ${inputClassName}`}
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
      />
      {label}
    </label>
  )
}

export default Checkbox
