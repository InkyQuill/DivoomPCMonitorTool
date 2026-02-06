import { describe, it, expect } from 'vitest'
import { formatPercent, formatBytes, formatTemperature, formatMetrics, getMetricsArray } from './formatters'
import type { SystemMetrics } from '@lib/types'

describe('formatPercent', () => {
	it('should format percentage with default 1 decimal', () => {
		expect(formatPercent(50.567)).toBe('50.6%')
	})

	it('should format percentage with custom decimals', () => {
		expect(formatPercent(50.567, 2)).toBe('50.57%')
		expect(formatPercent(50.567, 0)).toBe('51%')
	})

	it('should handle edge cases', () => {
		expect(formatPercent(0)).toBe('0.0%')
		expect(formatPercent(100)).toBe('100.0%')
		expect(formatPercent(99.99, 2)).toBe('99.99%')
	})
})

describe('formatBytes', () => {
	it('should format bytes in GB when >= 1 GB', () => {
		expect(formatBytes(5.432, 2)).toBe('5.43 GB')
		expect(formatBytes(1.0, 1)).toBe('1.0 GB')
	})

	it('should format bytes in MB when < 1 GB', () => {
		expect(formatBytes(0.5, 2)).toBe('512.00 MB')
		expect(formatBytes(0.001, 3)).toBe('1.024 MB')
	})

	it('should handle edge cases', () => {
		expect(formatBytes(0, 2)).toBe('0.00 GB')
		expect(formatBytes(0.999, 2)).toBe('1022.98 MB')
	})
})

describe('formatTemperature', () => {
	it('should format temperature in Celsius', () => {
		expect(formatTemperature(45.678)).toBe('45.7°C')
		expect(formatTemperature(0)).toBe('0.0°C')
	})

	it('should return N/A for null temperature', () => {
		expect(formatTemperature(null)).toBe('N/A')
	})
})

describe('formatMetrics', () => {
	const mockMetrics: SystemMetrics = {
		cpu: { usage_percent: 45.678, temperature: 50 },
		gpu: { usage_percent: 78.234, temperature: 65, vendor: 'NVIDIA' },
		memory: { usage_percent: 62.5, total_gb: 16, used_gb: 10, available_gb: 6 },
		storage: { usage_percent: 55.2, total_gb: 512, used_gb: 283, free_gb: 229, temperature: null },
		timestamp: Date.now(),
	}

	it('should format all metrics with percentages', () => {
		const result = formatMetrics(mockMetrics)
		expect(result).toBe('CPU: 45.7% | GPU: 78.2% | RAM: 62.5%')
	})

	it('should handle zero values', () => {
		const zeroMetrics: SystemMetrics = {
			...mockMetrics,
			cpu: { usage_percent: 0, temperature: null },
			gpu: { usage_percent: 0, temperature: null, vendor: 'AMD' },
			memory: { usage_percent: 0, total_gb: 16, used_gb: 0, available_gb: 16 },
		}
		const result = formatMetrics(zeroMetrics)
		expect(result).toBe('CPU: 0.0% | GPU: 0.0% | RAM: 0.0%')
	})
})

describe('getMetricsArray', () => {
	const mockMetrics: SystemMetrics = {
		cpu: { usage_percent: 45.678, temperature: 50 },
		gpu: { usage_percent: 78.234, temperature: 65, vendor: 'NVIDIA' },
		memory: { usage_percent: 62.5, total_gb: 16, used_gb: 10, available_gb: 6 },
		storage: { usage_percent: 55.2, total_gb: 512, used_gb: 283, free_gb: 229, temperature: null },
		timestamp: Date.now(),
	}

	it('should return array of formatted metric values', () => {
		const result = getMetricsArray(mockMetrics)
		expect(result).toEqual(['46', '78', '50', '65', '63', '55'])
	})

	it('should handle null temperatures by returning "0"', () => {
		const nullTempMetrics: SystemMetrics = {
			...mockMetrics,
			cpu: { usage_percent: 45, temperature: null },
			gpu: { usage_percent: 78, temperature: null, vendor: 'AMD' },
		}
		const result = getMetricsArray(nullTempMetrics)
		expect(result).toEqual(['45', '78', '0', '0', '63', '55'])
	})

	it('should round all values properly', () => {
		const result = getMetricsArray(mockMetrics)
		// CPU: 45.678 -> 46
		// GPU: 78.234 -> 78
		// Memory: 62.5 -> 63
		// Storage: 55.2 -> 55
		expect(result[0]).toBe('46') // CPU rounds up
		expect(result[1]).toBe('78') // GPU rounds down
		expect(result[4]).toBe('63') // Memory rounds up
		expect(result[5]).toBe('55') // Storage rounds down
	})
})
