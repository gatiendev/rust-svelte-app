<script lang="ts">
    import { auth } from "$lib/stores/auth";
    import { goto } from "$app/navigation";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";

    let loggingOut = false;

    // Redirect to login if not authenticated (and not loading)
    $: if (!$auth.loading && !$auth.user) {
        goto("/auth/login");
    }

    async function handleLogout() {
        loggingOut = true;
        await auth.logout();
        loggingOut = false;
        goto("/");
    }
</script>

{#if $auth.loading}
    <div class="container mx-auto px-4 py-16 flex justify-center">
        <div class="text-center text-gray-600 dark:text-gray-400">
            Loading profile...
        </div>
    </div>
{:else if $auth.user}
    <div class="container mx-auto px-4 py-16">
        <Card class="max-w-2xl mx-auto p-8">
            <h1 class="text-3xl font-bold mb-6 text-gray-900 dark:text-white">
                Profile
            </h1>
            <div class="space-y-4">
                <div>
                    <span class="font-semibold text-gray-700 dark:text-gray-300"
                        >User ID:</span
                    >
                    <span
                        class="ml-2 text-gray-900 dark:text-white font-mono text-sm"
                        >{$auth.user.id}</span
                    >
                </div>
                <div>
                    <span class="font-semibold text-gray-700 dark:text-gray-300"
                        >Username:</span
                    >
                    <span class="ml-2 text-gray-900 dark:text-white"
                        >{$auth.user.username}</span
                    >
                </div>
                <div>
                    <span class="font-semibold text-gray-700 dark:text-gray-300"
                        >Member since:</span
                    >
                    <span class="ml-2 text-gray-900 dark:text-white">
                        {new Date($auth.user.created_at).toLocaleDateString(
                            undefined,
                            {
                                year: "numeric",
                                month: "long",
                                day: "numeric",
                            },
                        )}
                    </span>
                </div>
                <div class="pt-6">
                    <Button
                        on:click={handleLogout}
                        disabled={loggingOut}
                        variant="outline"
                    >
                        {loggingOut ? "Logging out..." : "Logout"}
                    </Button>
                </div>
            </div>
        </Card>
    </div>
{/if}
