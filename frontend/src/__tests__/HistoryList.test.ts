import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import HistoryList from '../components/HistoryList.vue'

describe('HistoryList', () => {
  it('shows empty message when no history', () => {
    const wrapper = mount(HistoryList, {
      props: { history: [] }
    })
    expect(wrapper.find('.empty').text()).toBe('暂无历史记录')
  })

  it('renders history items with totals', () => {
    const wrapper = mount(HistoryList, {
      props: {
        history: [
          {
            taskId: 'abc12345',
            title: '测试报告',
            repoUrl: 'git@example.com:test/repo.git',
            sourceRef: 'feature',
            targetRef: 'main',
            createdAt: '2026-03-26 10:00:00',
            reportDir: '/tmp/report',
            htmlPath: '/tmp/report/report.html',
            markdownPath: '/tmp/report/report.md',
            jsonPath: '/tmp/report/report.json',
            totalIssues: 5,
            summary: {
              highRisk: 1,
              severe: 1,
              general: 2,
              suggestion: 1,
              high: 2,
              medium: 2,
              low: 1,
              total: 5
            }
          }
        ]
      }
    })
    expect(wrapper.find('.task-id').text()).toBe('abc12345')
    expect(wrapper.text()).toContain('总计 5')
    expect(wrapper.text()).toContain('高 2')
  })
})
