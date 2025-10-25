import React, { useState } from 'react'
import * as Yup from 'yup'
import { Formik, Form, Field, ErrorMessage, FormikHelpers } from 'formik'
import { useNavigate } from 'react-router-dom'

import { signIn, signUp } from '@behavior'
import { useAppDispatch } from '@hooks'
import { SignInReq, SignUpReq } from '@types'

const SignInSchema = Yup.object().shape({
  identifier: Yup.string().required('Username or email required'),
  password: Yup.string().required('Password required'),
})

const SignUpSchema = Yup.object().shape({
  username: Yup.string().required('Username required'),
  email: Yup.string().email('Invalid email').required('Email required'),
  password: Yup.string().required('Password required'),
  confirmPassword: Yup.string().oneOf(
    [Yup.ref('password')],
    'Passwords must match',
  ),
})

interface SignInValues {
  identifier: string
  password: string
}

const signInValues: SignInValues = {
  identifier: '',
  password: '',
}

interface SignUpValues {
  username: string
  email: string
  password: string
  confirmPassword: string
}

const signUpValues: SignUpValues = {
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
}

const AuthForm: React.FC = () => {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [isSignIn, setIsSignIn] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const onChangeSignIn = () => {
    setIsSignIn(!isSignIn)
    setError(null)
  }

  const onSubmit = async (
    values: SignInValues | SignUpValues,
    { setSubmitting, resetForm }: FormikHelpers<SignInValues | SignUpValues>,
  ) => {
    let user = null
    if (isSignIn) {
      user = await signIn(dispatch, values as SignInReq)
    } else {
      user = await signUp(dispatch, values as SignUpReq)
    }

    setSubmitting(false)
    resetForm()

    if (!user) {
      setError('Something went wrong! Try again soon.')
      return
    }

    navigate('/chess')
  }

  return (
    <div className='border-2 border-[#333] rounded-lg w-[300px] h-[500px] p-5'>
      <Formik
        initialValues={isSignIn ? signInValues : signUpValues}
        validationSchema={isSignIn ? SignInSchema : SignUpSchema}
        onSubmit={onSubmit}
      >
        {({ isSubmitting }) => (
          <div className={isSubmitting ? 'animate-pulse' : undefined}>
            <h2 className='font-semibold text-center mb-5 text-xl'>
              {isSignIn ? 'Sign In' : 'Sign Up'}
            </h2>
            <Form>
              <FormInput
                name={isSignIn ? 'identifier' : 'username'}
                label={isSignIn ? 'Username or Email' : 'Username'}
                type='text'
              />
              {!isSignIn && (
                <FormInput name='email' label='Email' type='email' />
              )}
              <FormInput name='password' label='Password' type='password' />
              {!isSignIn && (
                <FormInput
                  name='confirmPassword'
                  label='Confirm Password'
                  type='password'
                />
              )}
              <div className='flex justify-center'>
                <button type='submit' disabled={isSubmitting}>
                  Let&apos;s Go!
                </button>
              </div>
              {error && (
                <div className='mt-5'>
                  <FormError message={error} />
                </div>
              )}
              <p className='mt-5'>
                {isSignIn ? 'Need an account?' : 'Already have an account?'}
                <button
                  type='button'
                  onClick={onChangeSignIn}
                  className='text-blue-500 bg-inherit w-auto hover:border-[#242424] p-0 ml-1'
                >
                  {isSignIn ? 'Sign Up' : 'Sign In'}
                </button>
              </p>
            </Form>
          </div>
        )}
      </Formik>
    </div>
  )
}

interface FormInputProps {
  name: string
  label: string
  type: string
}

const FormInput: React.FC<FormInputProps> = ({ name, label, type }) => (
  <div className='flex flex-col mb-5'>
    <label htmlFor={name}>{label}</label>
    <Field type={type} name={name} id={name} className='rounded-lg p-1' />
    <ErrorMessage name={name} render={(msg) => <FormError message={msg} />} />
  </div>
)

interface FormErrorProps {
  message: string
}

const FormError: React.FC<FormErrorProps> = ({ message }) => (
  <div className='text-red-500 font-semibold text-sm'>{message}</div>
)

export default AuthForm
