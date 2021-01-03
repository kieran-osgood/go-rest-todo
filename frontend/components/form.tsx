import * as React from 'react';
import { useMutation } from 'react-query';
import axios from 'axios';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';

import styles from '../styles/Home.module.scss';

const todoForm = z.object({
  text: z.string().nonempty(),
});

type FormData = z.infer<typeof todoForm>;

const postTodo = async (formData: FormData) => {
  const { data } = await axios.post(
    '/api/todos', formData,
  );
  return data;
};

export default function Form() {
  const { register, handleSubmit, reset } = useForm({
    resolver: zodResolver(todoForm),
  });

  const { mutate } = useMutation((data: FormData) => postTodo(data), { onMutate: () => reset() });

  const submit = (data: FormData) => {
    mutate(data);
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit(submit)}>
      <label htmlFor="todo" className={styles.label}>
        <input type="text" id="todo" name="text" className={styles.input} ref={register} />
      </label>
      <button type="submit" className={styles.button}>Add</button>
    </form>
  );
}
