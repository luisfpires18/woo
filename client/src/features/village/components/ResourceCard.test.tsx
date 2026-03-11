import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { ResourceCard } from './ResourceCard';

// Mock the GameIcon component since it depends on the asset store
vi.mock('../../../components/GameIcon/GameIcon', () => ({
  GameIcon: ({ fallback }: { fallback: string }) => <span data-testid="game-icon">{fallback}</span>,
}));

describe('ResourceCard', () => {
  const defaultProps = {
    assetId: 'food',
    fallbackIcon: '🌾',
    label: 'Food',
    current: 450,
    max: 1200,
    rate: 8,
  };

  it('renders resource name, amount, and max', () => {
    render(<ResourceCard {...defaultProps} />);

    expect(screen.getByText('Food')).toBeInTheDocument();
    expect(screen.getByText('450')).toBeInTheDocument();
    expect(screen.getByText('1,200')).toBeInTheDocument();
  });

  it('renders positive rate with + prefix', () => {
    render(<ResourceCard {...defaultProps} rate={8} />);
    expect(screen.getByText('+8/s')).toBeInTheDocument();
  });

  it('renders negative rate without + prefix', () => {
    render(<ResourceCard {...defaultProps} rate={-3} />);
    expect(screen.getByText('-3/s')).toBeInTheDocument();
  });

  it('shows progress bar at correct percentage', () => {
    const { container } = render(<ResourceCard {...defaultProps} current={600} max={1200} />);
    const fill = container.querySelector('[class*="progressFill"]');
    expect(fill).not.toBeNull();
    expect(fill!.getAttribute('style')).toContain('width: 50%');
  });

  it('shows danger class when storage is nearly full (>=95%)', () => {
    const { container } = render(<ResourceCard {...defaultProps} current={1180} max={1200} />);
    const fill = container.querySelector('[class*="progressFill"]');
    expect(fill!.className).toContain('danger');
  });

  it('shows warning class when storage is filling (>=80%, <95%)', () => {
    const { container } = render(<ResourceCard {...defaultProps} current={1000} max={1200} />);
    const fill = container.querySelector('[class*="progressFill"]');
    expect(fill!.className).toContain('warning');
  });

  it('renders fallback icon via GameIcon', () => {
    render(<ResourceCard {...defaultProps} />);
    expect(screen.getByTestId('game-icon')).toHaveTextContent('🌾');
  });
});
