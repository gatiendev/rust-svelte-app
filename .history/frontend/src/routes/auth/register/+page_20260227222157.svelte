<script lang="ts">
  import Card from '$lib/components/ui/Card.svelte';
  import RegisterForm from '$lib/components/auth/RegisterForm.svelte';
  import { auth } from '$lib/stores/auth';

  $: loading = $auth.loading;
  $: error = $auth.error;

  function handleRegister(event: CustomEvent<{ email: string; password: string; name?: string }>) {
    const { email, password, name } = event.detail;
    auth.register(email, password, name);
  }
</script>

<div class="container mx-auto px-4 py-16 flex justify-center items-center min-h-[calc(100vh-200px)]">
  <div class="w-full max-w-md">
    <Card class="p-8 shadow-xl border border-gray-200 dark:border-gray-700">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Create an Account</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-2">Get started with your free account</p>
      </div>
      <RegisterForm on:submit={handleRegister} {loading} {error} />
      <p class="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
        Already have an account?
        <a href="/auth/login" class="text-primary-600 dark:text-primary-400 hover:underline font-medium">Sign in</a>
      </p>
    </Card>
  </div>
</div>