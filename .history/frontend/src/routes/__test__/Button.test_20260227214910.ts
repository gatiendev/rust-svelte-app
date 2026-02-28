import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import Button from '$lib/components/ui/Button.svelte';

describe('Button', () => {
    it('renders with default props', () => {
        render(Button, { props: { children: 'Click me' } });
        const button = screen.getByRole('button', { name: /click me/i });
        expect(button).toBeInTheDocument();
        expect(button).toHaveClass('bg-blue-600'); // primary variant
    });

    it('handles click events', async () => {
        const user = userEvent.setup();
        const onClick = vi.fn();
        render(Button, { props: { children: 'Click', onclick: onClick } });
        const button = screen.getByRole('button');
        await user.click(button);
        expect(onClick).toHaveBeenCalledOnce();
    });

    it('applies disabled state', () => {
        render(Button, { props: { disabled: true, children: 'Disabled' } });
        const button = screen.getByRole('button');
        expect(button).toBeDisabled();
        expect(button).toHaveClass('opacity-50');
    });

    it('renders with outline variant', () => {
        render(Button, { props: { variant: 'outline', children: 'Outline' } });
        const button = screen.getByRole('button');
        expect(button).toHaveClass('border');
    });
});