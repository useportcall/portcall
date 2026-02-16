<#import "template.ftl" as layout>
<@layout.emailLayout>
<h1>Action required</h1>
<p>Your administrator has requested that you update your <strong style="color: #0a0a0a;">${realmName}</strong> account by completing the following:</p>
<ul style="background-color: #f5f5f5; padding: 16px 16px 16px 32px; border-radius: 8px; margin: 16px 0; border: 1px solid #e5e5e5;">
    <#list requiredActions as action>
    <li style="color: #0a0a0a; font-size: 14px; margin: 6px 0;">${msg("requiredAction.${action}")}</li>
    </#list>
</ul>
<p>Click the button below to complete these actions:</p>
<p style="text-align: center; padding: 8px 0;">
    <a href="${link}" class="button">Update account</a>
</p>
<p>Or copy and paste this link into your browser:</p>
<p class="link-box">
    <a href="${link}">${link}</a>
</p>
<p>This link expires in <strong style="color: #0a0a0a;">${linkExpirationFormatter(linkExpiration)}</strong>.</p>
<p class="hint">If you didn't expect this, please contact your administrator.</p>
</@layout.emailLayout>
