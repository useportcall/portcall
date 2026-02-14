<#import "template.ftl" as layout>
<@layout.emailLayout>
<h1>Action Required</h1>
<p>Your administrator has requested that you update your <strong>${realmName}</strong> account by performing the following action(s):</p>
<ul style="background-color: #f4f4f5; padding: 16px 16px 16px 32px; border-radius: 6px; margin: 16px 0;">
    <#list requiredActions as action>
    <li style="color: #18181b; margin: 8px 0;">${msg("requiredAction.${action}")}</li>
    </#list>
</ul>
<p>Click the button below to complete these actions:</p>
<p style="text-align: center;">
    <a href="${link}" class="button">Update Account</a>
</p>
<p>Or copy and paste this link into your browser:</p>
<p style="word-break: break-all; background-color: #f4f4f5; padding: 12px; border-radius: 6px; font-size: 14px;">
    <a href="${link}">${link}</a>
</p>
<p>This link will expire within <strong>${linkExpirationFormatter(linkExpiration)}</strong>.</p>
<p style="color: #71717a; font-size: 14px; margin-top: 24px;">If you did not expect this request, please contact your administrator.</p>
</@layout.emailLayout>
